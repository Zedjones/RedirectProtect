using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

using System.Linq;
using System.Threading.Tasks;
using System.Threading;
using System.Collections.Generic;
using System;

using MongoDB.Driver;

namespace RedirectProtect.Services
{
    public class DeletionService : IHostedService, IDisposable
    {
        private readonly ILogger<DeletionService> _logger;
        private readonly RedirectService _redirectService;
        private List<Task> _timerTasks;
        private Dictionary<string, (CancellationTokenSource, Task)> _taskMap;
        private Task _watchTask;
        public DeletionService(ILogger<DeletionService> logger, RedirectService redirectService)
        {
            _logger = logger;
            _redirectService = redirectService;
            _timerTasks = new List<Task>();
            _taskMap = new Dictionary<string, (CancellationTokenSource, Task)>();
        }
        public Task StartAsync(CancellationToken stopToken)
        {
            var redirs = _redirectService.GetRedirects();
            foreach (var redir in redirs)
            {
                if (!(redir.ExpirationTime is null))
                {
                    if (redir.ExpirationTime < DateTime.UtcNow)
                    {
                        _redirectService.DeleteRedirect(redir);
                        _logger.LogInformation("Deleted {0}", redir.Path);
                    }
                    else
                    {
                        var tokenSource = CancellationTokenSource.CreateLinkedTokenSource(stopToken);
                        var redirTask = HandleRedirect(redir, tokenSource.Token);
                        _taskMap[redir.Path] = (tokenSource, redirTask);
                        _logger.LogInformation($"Created redirect handler for {redir.Path}");
                    }
                }
            }
            _watchTask = WatchCollection(stopToken);
            return _watchTask;
        }
        public Task StopAsync(CancellationToken stopToken)
        {
            //TODO: Check if this works
            var deleteTasks = _taskMap.Values.Select(tuple => tuple.Item2);
            deleteTasks.Append(_watchTask);
            return Task.WhenAll(deleteTasks);
        }
        public async Task HandleRedirect(Database.Models.Redirect redir, CancellationToken stopToken)
        {
            var timeToWait = DateTime.UtcNow - redir.ExpirationTime;
            await Task.Delay(timeToWait.Value.Milliseconds, stopToken);
            // Don't delete redirect if delay task was cancelled
            if (stopToken.IsCancellationRequested) return;
            _redirectService.DeleteRedirect(redir);
            _logger.LogInformation($"Deleted {redir.Path}");
        }
        public async Task WatchCollection(CancellationToken token = default)
        {
            using (var cursor = await _redirectService.GetRedirectCollection().WatchAsync(cancellationToken: token))
            {
                await cursor.ForEachAsync(change =>
                {
                    if (change.OperationType == ChangeStreamOperationType.Insert)
                    {
                        if (!(change.FullDocument.ExpirationTime is null))
                        {
                            if (change.FullDocument.ExpirationTime < DateTime.UtcNow)
                            {
                                _redirectService.DeleteRedirect(change.FullDocument);
                                _logger.LogInformation("Deleted {0}", change.FullDocument.Path);
                            }
                            else
                            {
                                var tokenSource = CancellationTokenSource.CreateLinkedTokenSource(token);
                                var redirTask = HandleRedirect(change.FullDocument, tokenSource.Token);
                                _taskMap[change.FullDocument.Path] = (tokenSource, redirTask);
                                _logger.LogInformation($"Created redirect handler for {change.FullDocument.Path}");
                            }
                        }
                    }
                    else if (change.OperationType == ChangeStreamOperationType.Delete)
                    {
                        //TODO: Take into account that we can't get full document from delete
                        if (!(change.FullDocument.ExpirationTime is null))
                        {
                            var (cancelSource, task) = _taskMap[change.FullDocument.Path];
                            if (!task.IsCompleted)
                            {
                                _logger.LogInformation($"Cancelling {change.FullDocument.Path}");
                                cancelSource.Cancel();
                            }
                        }
                    }
                }, cancellationToken: token);
            }
        }
        public void Dispose()
        {
            _timerTasks.ForEach(task => task.Dispose());
            _watchTask.Dispose();
        }
    }
}