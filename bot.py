from pyrogram import Client, filters
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import os

# Telegram Bot Configuration
API_ID = "your_api_id"
API_HASH = "your_api_hash"
BOT_TOKEN = "your_bot_token"

bot = Client("photoEnhanceBot", api_id=API_ID, api_hash=API_HASH, bot_token=BOT_TOKEN)

# Enhance Photo Function
def enhance_photo_with_automation(photo_path):
    # Selenium WebDriver Setup
    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument("--headless")  # For headless mode
    chrome_options.add_argument("--disable-gpu")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    
    driver = webdriver.Chrome(executable_path="path_to_chromedriver", options=chrome_options)
    driver.get("https://letsenhance.io/hi/boost")
    
    try:
        # Upload Photo
        upload_button = WebDriverWait(driver, 20).until(
            EC.presence_of_element_located((By.CSS_SELECTOR, "input[type='file']"))
        )
        upload_button.send_keys(photo_path)
        time.sleep(5)  # Wait for upload to complete

        # Start Enhancement
        process_button = WebDriverWait(driver, 20).until(
            EC.element_to_be_clickable((By.XPATH, "//button[contains(text(), 'Enhance')]"))
        )
        process_button.click()

        # Wait for Processing
        download_button = WebDriverWait(driver, 60).until(
            EC.element_to_be_clickable((By.XPATH, "//a[contains(@download, '.jpg')]"))
        )
        enhanced_photo_url = download_button.get_attribute("href")

        # Download Enhanced Photo
        enhanced_photo_path = "enhanced_photo.jpg"
        img_data = driver.get(enhanced_photo_url)
        with open(enhanced_photo_path, "wb") as f:
            f.write(img_data.content)
        driver.quit()
        return enhanced_photo_path
    except Exception as e:
        driver.quit()
        print(f"Error: {e}")
        return None

# Telegram Bot Handler for Photos
@bot.on_message(filters.photo & ~filters.edited)
async def handle_photo(client, message):
    photo = await message.download()  # Download photo sent by the user
    await message.reply("Processing your photo... Please wait.")

    # Enhance Photo
    enhanced_photo = enhance_photo_with_automation(photo)

    if enhanced_photo:
        await message.reply_document(enhanced_photo, caption="Here is your enhanced photo!")
        os.remove(enhanced_photo)  # Delete processed photo after sending
    else:
        await message.reply("Sorry, the enhancement process failed.")

# Run the Bot
if __name__ == "__main__":
    bot.run()
