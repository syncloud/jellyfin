from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By


def login(selenium, device_user, device_password):
    selenium.open_app()
    selenium.screenshot('index')
    selenium.click_by(By.XPATH, '//span[.="Next"]')
    selenium.wait_or_screenshot(EC.element_to_be_clickable((By.CSS_SELECTOR, "#txtManualName")))
    selenium.find_by_id("txtManualName").send_keys(device_user)
    password = selenium.find_by_id("txtManualPassword")
    password.send_keys(device_password)
    selenium.screenshot('login')
    password.send_keys(Keys.RETURN)
    selenium.screenshot('login_progress')
    selenium.find_by_css("span.material-icons.menu")
    selenium.screenshot('main')

