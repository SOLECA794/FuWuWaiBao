import httpx
import asyncio

async def main():
    url = 'http://127.0.0.1:8001/process'
    files = {'files': ('test_upload.txt', open('..\\data\\uploads\\test_upload.txt','rb'))}
    async with httpx.AsyncClient() as client:
        r = await client.post(url, files=files)
        print(r.status_code)
        print(r.text)

if __name__ == '__main__':
    asyncio.run(main())
