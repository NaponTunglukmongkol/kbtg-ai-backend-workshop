import { test, expect } from '@playwright/test';

test('Verify Contact information contains Nonthaburi', async ({ page }) => {
  // Navigate to the website
  await page.goto('https://www.kbtg.tech/th/home');

  // Check if the footer contains "ติดต่อ" (Thai for "Contact") and "Nonthaburi"
  const footerText = await page.textContent('footer');
  expect(footerText).toContain('ติดต่อ');
  expect(footerText).toContain('Nonthaburi');
});