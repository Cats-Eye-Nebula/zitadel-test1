module.exports = (on, config) => {
  // modify the config values
  config.defaultCommandTimeout = 20000

  //config.env.consoleUrl = "https://console.zitadel.ch"

  config.env.newEmail = "demo@caos.ch"
  config.env.newUserName = "demo"
  config.env.newFirstName = "demofirstname"
  config.env.newLastName = "demolastname"
  config.env.newPhonenumber = "+41 123456789"

  config.env.newMachineUserName = "demomachineusername"
  config.env.newMachineName = "demoname"
  config.env.newMachineDesription = "demodescription"


  // IMPORTANT return the updated config object
  return config

}
