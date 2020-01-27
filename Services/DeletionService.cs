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
        private Task _waitTask;
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
                    _logger.LogInformation("Deleted {0}", redir.Path);
                    _redirectService.DeleteRedirect(redir.Path);
                }
                else
                {
                }
            }
            return Task.WhenAll(_timerTasks);
        }
        public Task StopAsync(CancellationToken stopToken)
        {
            //TODO: Do something with cancellation token to correctly stop tasks
            return Task.CompletedTask;
        }
        public void Dispose()
        {
            _timerTasks.ForEach(task => task.Dispose());
            _waitTask.Dispose();
        }
    }
}