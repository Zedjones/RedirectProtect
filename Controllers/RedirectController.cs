using Microsoft.AspNetCore.Mvc;
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

        public RedirectController(RedirectService redirectService)
        {
            _redirectService = redirectService;
        }

        [HttpGet]
        public ActionResult<List<Redirect>> GetAction() => _redirectService.GetRedirects();

        [HttpPost]
        public ActionResult<Redirect> Create(RedirectDto redirect)
        {
            _redirectService.Create(redirect);

            return Ok();
        }
    }
}
