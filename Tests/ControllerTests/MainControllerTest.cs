using Xunit;
using Moq;
using System.Linq;
using RedirectProtect.Services;
using RedirectProtect.Database.Models;
using RedirectProtect.Controllers;

namespace RedirectProtect.Tests.ControllerTests
{
    public class MainControllerTest
    {
        [Theory]
        [InlineData("34672999")]
        [InlineData("40ce30aa")]
        [InlineData("1f9a3946")]
        public void GetRedirect_Missing_Path(string path)
        {
            var valid = new string[] {
                "9f2267a2",
                "065ce0c8",
                "76b27f02",
                "4a448e95",
                "c3d9c31e"
            };
            var actualRedirect = new Redirect { };
            var redirectServiceMock = new Mock<IRedirectService>();
            redirectServiceMock.Setup(mock =>
                mock.GetRedirect(It.IsAny<string>()))
                    .Returns((string s) => valid.Contains(s) ? actualRedirect : null);

            var mainController = new MainController(redirectServiceMock.Object);
            var res = mainController.GetRedirect(path);
            Assert.IsType<Microsoft.AspNetCore.Mvc.NotFoundResult>(res);
        }
        [Fact]
        public void GetRedirect_Correct_Path()
        {
            var actualRedirect = new Redirect { };
            var redirectServiceMock = new Mock<IRedirectService>();
            redirectServiceMock.Setup(mock => mock.GetRedirect("")).Returns(actualRedirect);
            var mainController = new MainController(redirectServiceMock.Object);
            var res = mainController.GetRedirect("");
            Assert.IsType<Microsoft.AspNetCore.Mvc.ViewResult>(res);
        }
        [Theory]
        [InlineData("34672999")]
        [InlineData("40ce30aa")]
        [InlineData("1f9a3946")]
        public void PostRedirect_Missing_Path(string path)
        {
            var valid = new string[] {
                "9f2267a2",
                "065ce0c8",
                "76b27f02",
                "4a448e95",
                "c3d9c31e"
            };
            var actualRedirect = new Redirect { };
            var redirectServiceMock = new Mock<IRedirectService>();
            redirectServiceMock.Setup(mock =>
                mock.GetRedirect(It.IsAny<string>()))
                    .Returns((string s) => valid.Contains(s) ? actualRedirect : null);

            var mainController = new MainController(redirectServiceMock.Object);
            var res = mainController.PostRedirect(path, "");
            Assert.IsType<Microsoft.AspNetCore.Mvc.NotFoundResult>(res);
        }
        [Fact]
        public void PostRedirect_Correct_Pass()
        {
            var hashedPassword = BCrypt.Net.BCrypt.HashPassword("aPassword");
            var actualRedirect = new Redirect
            {
                Password = hashedPassword,
                URL = "aurl.com"
            };
            var redirectServiceMock = new Mock<IRedirectService>();
            redirectServiceMock.Setup(mock => mock.GetRedirect("")).Returns(actualRedirect);
            var mainController = new MainController(redirectServiceMock.Object);
            var res = mainController.PostRedirect("", "aPassword");
            Assert.IsType<Microsoft.AspNetCore.Mvc.OkObjectResult>(res);
            Assert.Equal("aurl.com", ((Microsoft.AspNetCore.Mvc.OkObjectResult)res).Value);
        }
        [Theory]
        [InlineData("wrongPass")]
        [InlineData("ahAnotherwrong")]
        [InlineData("notrigh")]
        public void PostRedirect_Incorrect_Pass(string pass)
        {
            var hashedPassword = BCrypt.Net.BCrypt.HashPassword("aPassword");
            var actualRedirect = new Redirect
            {
                Password = hashedPassword,
                URL = "aurl.com"
            };
            var redirectServiceMock = new Mock<IRedirectService>();
            redirectServiceMock.Setup(mock => mock.GetRedirect("")).Returns(actualRedirect);
            var mainController = new MainController(redirectServiceMock.Object);
            var res = mainController.PostRedirect("", pass);
            Assert.IsType<Microsoft.AspNetCore.Mvc.UnauthorizedObjectResult>(res);
            Assert.Equal("Incorrect password provided",
                ((Microsoft.AspNetCore.Mvc.UnauthorizedObjectResult)res).Value);
        }
    }
}