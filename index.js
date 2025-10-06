import puppeteer from 'puppeteer';

async function getNotes (page) {
  const TOP_NOTES = 'Top Notes'
  const MIDDLE_NOTES = 'Middle Notes';
  const BASE_NOTES = 'Base Notes';

  const pyramid = await page.$('#pyramid');
  const text = await page.evaluate(el => el.innerText, pyramid);
  const tokens = text.split('\n').map(t => t.trim()).slice(4);

  const mnIdx = tokens.indexOf(MIDDLE_NOTES);
  const bnIdx = tokens.indexOf(BASE_NOTES);

  return {
    [TOP_NOTES]: tokens.slice(0, mnIdx),
    [MIDDLE_NOTES]: tokens.slice(mnIdx + 1, bnIdx),
    [BASE_NOTES]: tokens.slice(bnIdx + 1, tokens.length - 1)
  };
}

async function run () {
  const browser = await puppeteer.launch({
    headless: false
  });
  const page = await browser.newPage();

  const testPage = 'https://www.fragrantica.com/perfume/Lattafa-Perfumes/Khamrah-75805.html'; 
  await page.goto(testPage);

  const notes = await getNotes(page);

  console.log('-- results --');
  console.log(notes);
  console.log('-------------');

  browser.close();
}

run();