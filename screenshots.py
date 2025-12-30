#!/usr/bin/env python3
"""
Playwright script to capture screenshots of the Grafana Parquet S3 datasource plugin.
"""

import asyncio
from playwright.async_api import async_playwright


async def main():
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        context = await browser.new_context(viewport={"width": 1920, "height": 1080})
        page = await context.new_page()

        # Login to Grafana
        print("Logging into Grafana...")
        await page.goto("http://localhost:3000/login")
        await page.fill('input[name="user"]', "admin")
        await page.fill('input[name="password"]', "admin123")
        await page.click('button[type="submit"]')

        # Wait for login to complete
        await page.wait_for_url("**/", timeout=10000)
        print("Logged in successfully")

        # Screenshot 1: Data Sources list
        print("Capturing Data Sources list...")
        await page.goto("http://localhost:3000/connections/datasources")
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(2)  # Extra wait for UI to render
        await page.screenshot(path="screenshot_datasources_list.png", full_page=False)
        print("Saved: screenshot_datasources_list.png")

        # Screenshot 2: Datasource configuration page
        print("Capturing Datasource configuration...")
        # Click on the Parquet S3 datasource
        await page.click('text=Parquet S3')
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(2)
        await page.screenshot(path="screenshot_datasource_config.png", full_page=False)
        print("Saved: screenshot_datasource_config.png")

        # Screenshot 3: Explore view with query
        print("Capturing Explore view...")
        await page.goto("http://localhost:3000/explore")
        await page.wait_for_load_state("networkidle")
        await asyncio.sleep(2)

        # Select the Parquet S3 datasource if not already selected
        try:
            # Click on datasource picker
            datasource_picker = page.locator('[data-testid="data-source-picker"]').first
            if await datasource_picker.is_visible():
                await datasource_picker.click()
                await asyncio.sleep(1)
                # Select Parquet S3
                await page.click('text=Parquet S3')
                await asyncio.sleep(2)
        except Exception as e:
            print(f"Datasource picker: {e}")

        # Try to select bucket and run query
        try:
            # Wait for the query editor to load
            await asyncio.sleep(2)

            # Select bucket from dropdown
            bucket_select = page.locator('input[placeholder="Select bucket..."]').first
            if await bucket_select.is_visible():
                await bucket_select.click()
                await asyncio.sleep(1)
                await page.click('text=parquet-data')
                await asyncio.sleep(2)

            # Select time column
            time_select = page.locator('input[placeholder="(none)"]').first
            if await time_select.is_visible():
                await time_select.click()
                await asyncio.sleep(1)
                # Try to select timestamp
                timestamp_option = page.locator('text=timestamp').first
                if await timestamp_option.is_visible():
                    await timestamp_option.click()
                    await asyncio.sleep(2)

            # Run the query
            run_button = page.locator('button:has-text("Run query")').first
            if await run_button.is_visible():
                await run_button.click()
                await asyncio.sleep(3)  # Wait for query results

        except Exception as e:
            print(f"Query setup: {e}")

        await page.screenshot(path="screenshot_explore_view.png", full_page=False)
        print("Saved: screenshot_explore_view.png")

        # Screenshot 4: Explore view with data (full page)
        await page.screenshot(path="screenshot_explore_full.png", full_page=True)
        print("Saved: screenshot_explore_full.png")

        await browser.close()
        print("\nAll screenshots captured successfully!")


if __name__ == "__main__":
    asyncio.run(main())
