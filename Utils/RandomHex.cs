using System;
using System.Linq;

namespace RedirectProtect.Utils
{
    public class RandomHex
    {
        // From https://stackoverflow.com/questions/1054076/randomly-generated-hexadecimal-number-in-c-sharp
        static Random random = new Random();
        public static string GetRandomHexNumber(int digits)
        {
            var buffer = new byte[digits / 2];
            random.NextBytes(buffer);
            var result = String.Concat(buffer.Select(x => x.ToString("X2")).ToArray());
            if (digits % 2 == 0)
                return result;
            return result + random.Next(16).ToString("X");
        }
    }
}