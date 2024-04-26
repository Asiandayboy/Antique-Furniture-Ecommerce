import { Builder, By, Browser, until } from "selenium-webdriver";


function successPrint(msg) {
    console.log(`\x1b[32m${msg}\x1b[0m`)
}

function failPrint(msg) {
    console.log(`\x1b[31m${msg}\x1b[0m`)
}




async function testBlankField(testName, fieldName) {
    const driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/signup")


    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys(fieldName == "username" ? "" : "bruh")

    const emailInput = await driver.findElement({id: "email"})
    emailInput.sendKeys(fieldName == "email" ? "" : "123@gmail.com")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys(fieldName == "password" ? "" : "bruh")


    const confirmInput = await driver.findElement({id: "confirm"})
    confirmInput.sendKeys(fieldName == "confirm" ? "" : "bruh")
    

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    const errMsg = await driver.wait(until.elementLocated(By.css(".signup_err-msg")), 3000).getText()

    if (errMsg != "Fields cannot be blank") {
        failPrint(`Test ${testName}: failed -> expected: ${"Fields cannot be blank"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}

async function testUsedUsername(testName) {
    const driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/signup")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("bruh")

    const emailInput = await driver.findElement({id: "email"})
    emailInput.sendKeys("123@gmail.com")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("bruh")


    const confirmInput = await driver.findElement({id: "confirm"})
    confirmInput.sendKeys("bruh")
    

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    const errMsg = await driver.wait(until.elementLocated(By.css(".signup_err-msg")), 3000).getText()

    if (errMsg != "Username is taken") {
        failPrint(`Test ${testName}: failed -> expected: ${"Username is taken"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}

async function testUsedEmail(testName) {
    const driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/signup")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("bruhbro")

    const emailInput = await driver.findElement({id: "email"})
    emailInput.sendKeys("bruh@gmail.com")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("bruh")


    const confirmInput = await driver.findElement({id: "confirm"})
    confirmInput.sendKeys("bruh")
    

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    const errMsg = await driver.wait(until.elementLocated(By.css(".signup_err-msg")), 3000).getText()

    if (errMsg != "Email is taken") {
        failPrint(`Test ${testName}: failed -> expected: ${"Email is taken"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}

async function testMismatchedPassword(testName) {
    const driver = await new Builder().forBrowser(Browser.EDGE).build()

    await driver.get("http://127.0.0.1:5173/signup")

    const usernameInput = await driver.findElement({id: "username"})
    usernameInput.sendKeys("bruhbro")

    const emailInput = await driver.findElement({id: "email"})
    emailInput.sendKeys("bruhbro@gmail.com")

    const passwordInput = await driver.findElement({id: "password"})
    passwordInput.sendKeys("bruh")


    const confirmInput = await driver.findElement({id: "confirm"})
    confirmInput.sendKeys("bruh123")
    

    const submitButton = await driver.findElement({name: "submit"})
    await submitButton.click()

    const errMsg = await driver.wait(until.elementLocated(By.css(".signup_err-msg")), 3000).getText()

    if (errMsg != "Passwords do not match") {
        failPrint(`Test ${testName}: failed -> expected: ${"Passwords do not match"}, got: ${errMsg}`)
    } else {
        successPrint(`Test ${testName}: passed!`)
    }

    await driver.close()
}


testBlankField("1", "username")
testBlankField("2", "email")
testBlankField("3", "password")
testBlankField("4", "confirm")

testUsedUsername("5")
testUsedEmail("6")
testMismatchedPassword("7")
