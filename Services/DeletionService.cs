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
        private List<Timer> _timers;
        public DeletionService(ILogger<DeletionService> logger, RedirectService redirectService)
        {
            _logger = logger;
            _redirectService = redirectService;
            _timers = new List<Timer>();
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
                    _timers.Add(new Timer(
                        (state) => {
                            _redirectService.DeleteRedirect(redir.Path);
                            _logger.LogInformation("Deleted {0}", redir.Path);
                        }
                    ));
                }
            }
            return Task.CompletedTask;
        }
        public Task StopAsync(CancellationToken stopToken)
        {
            _timers.ForEach(timer => timer.Change(Timeout.Infinite, 0));

            return Task.CompletedTask;
        }
        public void Dispose()
        {
            _timers.ForEach(timer => timer.Dispose());
        }
    }
}