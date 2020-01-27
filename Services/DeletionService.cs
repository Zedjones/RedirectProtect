using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
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
        private Dictionary<string, CancellationToken> _taskMap;
        private Task _watchTask;
        public DeletionService(ILogger<DeletionService> logger, RedirectService redirectService)
        {
            _logger = logger;
            _redirectService = redirectService;
            _timerTasks = new List<Task>();
        }
        public Task StartAsync(CancellationToken stopToken)
        {
            var redirs = _redirectService.GetRedirects();
            foreach (var redir in redirs)
            {
                if (redir.ExpirationTime < DateTime.UtcNow)
                {
                    _redirectService.DeleteRedirect(redir);
                    _logger.LogInformation("Deleted {0}", redir.Path);
                }
                else
                {
                }
            }
            _watchTask = WatchCollection();
            return _watchTask;
        }
        public Task StopAsync(CancellationToken stopToken)
        {
            //TODO: Do something with cancellation token to correctly stop tasks
            return Task.CompletedTask;
        }
        public async Task HandleRedirect(Database.Models.Redirect redir, CancellationToken stopToken)
        {
            var timeToWait = DateTime.UtcNow - redir.ExpirationTime;
            await Task.Delay(timeToWait.Value.Milliseconds, stopToken);
            _redirectService.DeleteRedirect(redir);
            _logger.LogInformation("Deleted {0}", redir.Path);
        }
        public async Task WatchCollection(CancellationToken token = default)
        {
            using (var cursor = await _redirectService.GetRedirectCollection().WatchAsync(cancellationToken: token))
            {
                await cursor.ForEachAsync(change =>
                {
                    if (change.OperationType == ChangeStreamOperationType.Insert)
                    {
                        
                    }
                    else if (change.OperationType == ChangeStreamOperationType.Delete)
                    {

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