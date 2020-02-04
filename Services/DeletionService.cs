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
        private Dictionary<string, (CancellationTokenSource, Task)> _taskMap;
        private List<Task> _taskList;
        private Task _waitTask;
        private CancellationToken mainToken;
        public DeletionService(ILogger<DeletionService> logger, RedirectService redirectService)
        {
            _logger = logger;
            _redirectService = redirectService;
            _taskList = new List<Task>();
            _taskMap = new Dictionary<string, (CancellationTokenSource, Task)>();
        }
        public Task StartAsync(CancellationToken stopToken)
        {
            var redirs = _redirectService.GetRedirects();
            mainToken = stopToken;
            foreach (var redir in redirs)
            {
                ProcessRedirect(redir);
            }
            _waitTask = Task.WhenAll(_taskList);
            return _waitTask;
        }
        public Task StopAsync(CancellationToken stopToken)
        {
            _logger.LogInformation("Stopping deletion service");
            foreach (var tokenSource in _taskMap.Values.Select(tuples => tuples.Item1))
            {
                tokenSource.Cancel();
            }
            return _waitTask;
        }
        public async Task HandleRedirect(Database.Models.Redirect redir, CancellationToken stopToken)
        {
            var timeToWait = redir.ExpirationTime - DateTime.UtcNow;
            await Task.Delay(timeToWait.Value, stopToken);
            // Don't delete redirect if delay task was cancelled or if it was manually deleted
            _logger.LogInformation($"Is cancelled: {stopToken.IsCancellationRequested}");
            if (stopToken.IsCancellationRequested || !(_redirectService.RedirectExists(redir))) return;
            _redirectService.DeleteRedirect(redir);
            _logger.LogInformation($"Deleted {redir.Path}");
        }
        public void ProcessRedirect(Database.Models.Redirect redir)
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
                    var tokenSource = CancellationTokenSource.CreateLinkedTokenSource(mainToken);
                    var redirTask = HandleRedirect(redir, tokenSource.Token);
                    _taskMap[redir.Path] = (tokenSource, redirTask);
                    _logger.LogInformation($"Created redirect handler for {redir.Path}");
                    _taskList.Add(redirTask);
                }
            }
        }
        public void Dispose()
        {
            _logger.LogInformation("Disposing of tasks");
            _taskList.ForEach(task => task.Dispose());
            _waitTask.Dispose();
        }
    }
}