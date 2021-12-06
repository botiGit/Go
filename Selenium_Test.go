package selenium_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

//	Programa de prueba que va a este git y lista todos los repositorios que tenga
//
//  If you want to actually run this example:
//
//   1. Ensure the file paths at the top of the function are correct.
//   2. Remove the word "Example" from the comment at the bottom of the
//      function.
//   3. Run:
//      go test -test.run=Example$ github.com/tebeka/selenium


func Example() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		// These paths will be different on your system.
		seleniumPath     = "vendor/selenium-server-standalone-3.4.jar"
		chromeDriverPath = "vendor/geckodriver-v0.18.0-linux64"
		port             = 8080
	)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(geckoDriverPath),
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()


	if err := wd.Get("https://github.com/botiGit"); err != nil {
		panic(err)
	}

	//Buscamos el botón de repositorios y clickamos
	reposBtn, err := wd.FindElement(selenium.ByXPATH, "//*[@id="js-pjax-container"]/div[1]/div/div/div[2]/div/nav/a[2]")
	if err != nil {
		panic(err)
	}
	if err := reposBtn.Click(); err != nil {
		panic(err)
	}

	//Array con todos los h3s en repositories
	repositories, err := wd.FindElements(selenium.ByTagName, "h3")
	if err != nil {
		panic(err)
	}

	var hijos [] WebElement
	j := 0
	//Quitamos los dos primeros y buscamos dentro de cada uno de los h3 el tag "a" que tiene el título del repo
	for i := 2; i < len(repositories); i++ {
    	hijos[j] = repositories[i].FindElement(selenium.ByTagName, "a")
    	j++
    	if err != nil {
			panic(err)
		}
	}

	for i := 0; i < len(hijos); i++ {
    	fmt.Println(hijos[i].Text())
	}
	

}