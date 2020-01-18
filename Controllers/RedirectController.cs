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

        [HttpGet("{id:length(24)}", Name = "GetRedirect")]
        public ActionResult<Redirect> GetAction(string id)
        {
            var redirect = _redirectService.GetRedirect(id);

            if (redirect is null)
            {
                return NotFound();
            }

            return redirect;
        }

        [HttpPost]
        public ActionResult<Redirect> Create(Redirect redirect)
        {
            _redirectService.Create(redirect);

            return CreatedAtRoute("GetBook", new { id = redirect.Id });
        }
    }
}
