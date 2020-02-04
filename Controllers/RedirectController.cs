using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Hosting;
using System.Net;
using System.Linq;
using RedirectProtect.Services;
using RedirectProtect.Database.Models;
using System.Collections.Generic;

namespace RedirectProtect.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class RedirectController : ControllerBase
    {
        private readonly RedirectService _redirectService;
        private readonly ILogger<RedirectController> _logger;
        private readonly DeletionService _deletionService;

        public RedirectController(RedirectService redirectService, ILogger<RedirectController> logger,
                                  IEnumerable<IHostedService> services)
        {
            _redirectService = redirectService;
            _logger = logger;
            // Have to do this because we can't directly injected DeletionService
            _deletionService = (DeletionService)services.FirstOrDefault(w => w.GetType() == typeof(DeletionService));
        }

        [HttpGet]
        public ActionResult<List<Redirect>> GetAction() => _redirectService.GetRedirects();

        [HttpPost]
        public ActionResult<Redirect> Create(RedirectDto redirect)
        {
            try
            {
                var newRedir = _redirectService.Create(redirect);
                _deletionService.ProcessRedirect(newRedir);
            }
            catch (System.TimeoutException te)
            {
                _logger.LogError(te.Message);
                return Problem(detail: "Could not connect to MongoDB client",
                               statusCode: (int)HttpStatusCode.InternalServerError);
            }

            return Ok();
        }
    }
}
