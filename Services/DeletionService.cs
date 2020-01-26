using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
using System.Threading;
using System.Collections.Generic;
using System;

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
        }
        public Task StartAsync(CancellationToken stopToken)
        {
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