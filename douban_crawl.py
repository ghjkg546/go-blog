from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.chrome.options import Options
import time
import requests
import os
import time
import requests
from PIL import Image
import io
import os


import requests
import os

def postToServer(id,content,img):
# 设置目标接口 URL
    url = "https://api.shareziyuan.email/adminapi/info_upload"  # 替换为实际接口地址

    # 要发送的数据
    data = {
        "id": id,
        "content": content
    }

    # 图片路径
    image_path = os.path.join(os.getcwd(), img)

    # 检查图片是否存在
    if not os.path.exists(image_path):
        raise FileNotFoundError(f"图片文件不存在: {image_path}")

    # 准备图片文件上传
    files = {
        "img": open(image_path, "rb")
    }

    # try:
        # 发送 POST 请求
    response = requests.post(url, data=data, files=files)
    print(response.text)
  
    files["img"].close()


# 配置 Selenium WebDriver
chrome_options = Options()
chrome_options.add_argument("--headless")  # 如果需要无头模式运行，可以取消注释
chrome_options.add_argument("--disable-gpu")
chrome_options.add_argument("--no-sandbox")

# 替换为你的 ChromeDriver 路径
service = Service("")  # 请替换为 chromedriver 的实际路径
driver = webdriver.Chrome(service=service, options=chrome_options)


def crawl_info(id,doubanId):

    # 目标 URL
    url = f"https://movie.douban.com/subject/{doubanId}/"

    try:
        # 打开目标页面
        driver.get(url)
        time.sleep(2)  # 等待页面加载，具体时间根据需要调整

        # 获取 class="subject clearfix" 的内容
        subject_element = driver.find_element(By.CLASS_NAME, "subject.clearfix")
        if subject_element:
            print("Subject内容:", subject_element.get_attribute("outerHTML"))

        # 获取 id="link-report-intra" 的内容
        link_report_element = driver.find_element(By.ID, "link-report-intra")
        if link_report_element:
            print("Brief内容:", link_report_element.get_attribute("outerHTML"))

        # 获取 #mainpic 下 img 的 src
        mainpic_img = driver.find_element(By.CSS_SELECTOR, "#mainpic img")
        if mainpic_img:
            img_url = mainpic_img.get_attribute("src")
            print("图片链接:", img_url)
            print("图片链接:", img_url)

            # 下载图片
            img_response = requests.get(img_url)
            img_response.raise_for_status()

            # 转换为 JPG 并保存
            img = Image.open(io.BytesIO(img_response.content))
            img_name = os.path.splitext(os.path.basename(img_url))[0] + ".jpg"
            img = img.convert("RGB")  # 转换为 RGB 模式，适配 JPG 格式
            img.save(img_name, "JPEG")
            print(f"图片已保存为 {img_name}")

        
        postToServer(id,subject_element.get_attribute("outerHTML")+link_report_element.get_attribute("outerHTML"),img_name)

    except Exception as e:
        print(f"发生错误: {e}")

       

crawl_info(2734,2158647)
crawl_info(2735,30377729)

 # 关闭浏览器
driver.quit()