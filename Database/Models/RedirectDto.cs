using System;

namespace RedirectProtect.Database.Models
{
    public class RedirectDto
    {
        public string URL { get; set; }
        public string Password { get; set; }
        public DateTime TTL { get; set; }
    }
}