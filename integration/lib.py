from selenium.webdriver.common.keys import Keys


def login(selenium, device_user, device_password):
    selenium.open_app()
    selenium.screenshot('index')
    selenium.find_by_id("txtManualName").send_keys(device_user)
    password = selenium.find_by_id("txtManualPassword")
    password.send_keys(device_password)
    selenium.screenshot('login')
    password.send_keys(Keys.RETURN)
    selenium.screenshot('login_progress')
    selenium.find_by_xpath("//button[@data-qa='sidebar-search']")
    selenium.screenshot('main')

