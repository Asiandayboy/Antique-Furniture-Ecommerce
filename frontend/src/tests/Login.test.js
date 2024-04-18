import { Builder, By, Browser, until } from "selenium-webdriver";

    // const textInput = await driver.findElement({id: "my-text-id"})
    // textInput.sendKeys("Hello World!")

    // const passwordInput = await driver.findElement({name: "my-password"})
    // passwordInput.sendKeys("password123bruh")
    
    // const textArea = await driver.findElement({name: "my-textarea"})
    // textArea.sendKeys("some random text\n")

    // const colorPicker = await driver.findElement({name: "my-colors"})
    // colorPicker.sendKeys("#FF6347")

    // const readonlyInput = await driver.findElement({name: "my-readonly"})
    // const readonlyText = await readonlyInput.getAttribute("value")
    // textArea.sendKeys(readonlyText)
    // console.log("READONLY TEXT:", readonlyText)

    // const datePicker = await driver.findElement({name: "my-date"})
    // datePicker.sendKeys("03/20/2024")

    // const dropdownSelect = await driver.findElement({name: "my-select"})
    // dropdownSelect.sendKeys("One")

    // const dropdownList = await driver.findElement({name: "my-datalist"})
    // dropdownList.sendKeys("Seattle")

    // const checkedCheckbox = await driver.findElement({id: "my-check-1"})
    // checkedCheckbox.click()

    // const defaultCheckbox = await driver.findElement({id: "my-check-2"})
    // defaultCheckbox.click()

    // const checkedRadio = await driver.findElement({id: "my-radio-2"})
    // checkedRadio.click()

    // const defaultRadio = await driver.findElement({id: "my-radio-2"})
    // defaultRadio.click()
    
    // setTimeout(async () => {
    //     const submitButton = await driver.findElement({className: "btn btn-outline-primary mt-3"})
    //     submitButton.click()

    // }, 5000)


function successPrint(msg) {
    console.log(`\x1b[32m${msg}\x1b[0m`)
}

function failPrint(msg) {
    console.log(`\x1b[31m${msg}\x1b[0m`)
}

async function testBlankPassword(testName) {
    let driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/login")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("bruh")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("")

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()


    const errMsg = await driver.wait(until.elementLocated(By.css(".login_err-msg")), 3000).getText()

    if (errMsg != "Password cannot be blank") {
        failPrint(`Test ${testName}: failed -> expected: ${"Password cannot be blank"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}

async function testBlankUsername(testName) {
    let driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/login")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("password123")

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()


    const errMsg = await driver.wait(until.elementLocated(By.css(".login_err-msg")), 3000).getText()

    if (errMsg != "Username cannot be blank") {
        failPrint(`Test ${testName}: failed -> expected: ${"Username cannot be blank"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}

// test must be run with server running
async function testValidLogin(testName) {
    let driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/login")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("varus")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("varuslover")

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    // if the client gets redirected to the dashboard, then that means login was successful
    await driver.wait(until.urlContains("dashboard"), 3000)

    const title = await driver.findElement({tagName: "h1"}).getText()

    if (title == "Dashboard") {
        successPrint(`Test ${testName}: passed!`)
    } else {
        failPrint("Test ${testName}: failed -> expected a valid login, got an invalid login")
    }

    await driver.close()
}

// test must be run with server running
async function testInvalidLogin(testName) {
    let driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/login")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("bruh")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("bruh")

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    const errMsg = await driver.wait(until.elementLocated(By.css(".login_err-msg")), 3000).getText()

    if (errMsg != "Invalid login") {
        failPrint(`Test ${testName}: failed -> expected: ${"Invalid login"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}



testBlankUsername("1")
testBlankPassword("2")
testValidLogin("3")
testInvalidLogin("4")