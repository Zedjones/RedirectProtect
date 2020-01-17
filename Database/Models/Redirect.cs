namespace RedirectProtect.Database.Models
{
    public class Redirect
    {
        public string Path { get; set; }
        public string URL { get; set; }
        public string Password { get; set; }
        public string TTL { get; set; }
    }
}