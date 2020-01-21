using Microsoft.AspNetCore.Mvc;
using RedirectProtect.Services;

namespace RedirectProtect.Controllers
{
    [ApiController]
    [Route("/")]
    public class MainController : ControllerBase
    {
        private readonly RedirectService _redirectService;
        public MainController(RedirectService redirectService)
        {
            _redirectService = redirectService;
        }
        [HttpGet("{path:length(8)}", Name = "GetRedirect")]
        public IActionResult GetRedirect(string path)
        {
            var redirect = _redirectService.GetRedirect(path);

            if (redirect is null)
            {
                return NotFound();
            }

            return Ok(redirect.Path);
        }
    }
}