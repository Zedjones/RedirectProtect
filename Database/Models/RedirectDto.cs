using System;
using System.ComponentModel;
using Newtonsoft.Json;

namespace RedirectProtect.Database.Models
{
    public class RedirectDto
    {
        public string URL { get; set; }
        public string Password { get; set; }
        [DefaultValue(null)]
        [JsonProperty(DefaultValueHandling = DefaultValueHandling.Populate)]
        public DateTime? TTL { get; set; }
    }
}