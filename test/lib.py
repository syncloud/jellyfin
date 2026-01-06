from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By


def login(selenium, device_user, device_password, mode):
    selenium.open_app()
    selenium.screenshot(mode+'-index')
    # selenium.click_by(By.XPATH, '//span[.="Next"]')
    selenium.wait_or_screenshot(EC.element_to_be_clickable((By.CSS_SELECTOR, "#txtManualName")))
    selenium.find_by_id("txtManualName").send_keys(device_user)
    password = selenium.find_by_id("txtManualPassword")
    password.send_keys(device_password)
    selenium.screenshot(mode+'-login')
    password.send_keys(Keys.RETURN)
    selenium.screenshot(mode+'-login_progress')
    selenium.find_by(By.XPATH, "//h2[.='Nothing here.']")
    selenium.screenshot(mode+'-main')

def scan(selenium, mode):
    selenium.click_by(By.XPATH, "//button[@title='Menu']")
    selenium.click_by(By.XPATH, "//span[.='Dashboard']")
    selenium.click_by(By.XPATH, "//span[.='Scan All Libraries']")
    selenium.screenshot(mode+'-scan')
    selenium.find_by(By.XPATH, "//span[.='Running Tasks']")
    selenium.invisible_by(By.XPATH, "//span[.='Running Tasks']")
    selenium.screenshot(mode+'-scan-done')
