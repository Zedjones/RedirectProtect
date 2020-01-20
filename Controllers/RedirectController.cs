using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using System.Net;
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

        public RedirectController(RedirectService redirectService, ILogger<RedirectController> logger)
        {
            _redirectService = redirectService;
            _logger = logger;
        }

        [HttpGet]
        public ActionResult<List<Redirect>> GetAction() => _redirectService.GetRedirects();

        [HttpPost]
        public ActionResult<Redirect> Create(RedirectDto redirect)
        {
            try
            {
                _redirectService.Create(redirect);
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
