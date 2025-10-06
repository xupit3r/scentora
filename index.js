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

async function getAccords (page) {
  const accords = await page.$$eval('.accord-bar', results => {
    return results.map(el => {
      const accord = el.innerText.trim();
      const parentWidth = el.parentElement.offsetWidth;
      const myWidth = el.offsetWidth;
      const percent = Math.round(myWidth / parentWidth * 100);

      return {
        accord,
        percent
      }
    });
  });

  return accords;
}

async function getDescription (page) {
  return await page.$eval('.fragrantica-blockquote', q => q.innerText.trim());
}

async function getPageInfo (url) {
  const browser = await puppeteer.launch({
    headless: false
  });
  const page = await browser.newPage();
 
  await page.goto(url);

  const notes = await getNotes(page);
  const accords = await getAccords(page);
  const description = await getDescription(page);

  console.log('-- results --');
  console.log(notes);
  console.log(accords);
  console.log(description);
  console.log('-------------');

  browser.close();
}

getPageInfo('https://www.fragrantica.com/perfume/Lattafa-Perfumes/Khamrah-75805.html');