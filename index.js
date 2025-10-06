import puppeteer from 'puppeteer';

async function getNotes (page) {
  const pyramid = await page.$('#pyramid');
  const text = await page.evaluate(el => el.innerText, pyramid);
  const tokens = text.split('\n').map(t => t.trim()).slice(4);

  const mnIdx = tokens.indexOf('Middle Notes');
  const bnIdx = tokens.indexOf('Base Notes');

  return {
    "Top Notes": tokens.slice(0, mnIdx),
    "Middle Notes": tokens.slice(mnIdx + 1, bnIdx),
    "Base Notes": tokens.slice(bnIdx + 1, tokens.length - 1)
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