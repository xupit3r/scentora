import puppeteer from 'puppeteer';

async function getFragranceNotes (page) {
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

async function getFrangranceAccords (page) {
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

async function getFrangranceDescription (page) {
  return await page.$eval('.fragrantica-blockquote', q => q.innerText.trim());
}

async function getPageInfo (browser, url) {
  const page = await browser.newPage();
 
  await page.goto(url);

  const notes = await getFragranceNotes(page);
  const accords = await getFrangranceAccords(page);
  const description = await getFrangranceDescription(page);

  console.log('-- results --');
  console.log(notes);
  console.log(accords);
  console.log(description);
  console.log('-------------');
}

async function collectLinks (page) {
  try {
    for (let i = 0; i < 30; i++) {
      await loadMoreResults(page);
      await delay(500);
    }
  } catch (err) {
    console.error(err);
  }
  
  return page.$$eval('.cell.card.fr-news-box a', results => {
    return results.map(el => el.href);
  });
}

async function loadMoreResults (page) {
  await page.locator('text/Show more results')
            .setEnsureElementIsInTheViewport(true)
            .click();
}

function delay (time) {
  return new Promise((resolve) => {
    setTimeout(resolve, time);
  })
}

async function getNotes (page) {
  await page.goto("https://www.fragrantica.com/notes/");

  return page.$$eval('.notebox a', results => {
    return results.map(el => {
      return {
        note: el.innerText.trim(),
        url: el.href
      }
    });
  }); 
}

async function run () {
  const browser = await puppeteer.launch({
    headless: false
  });
  const page = await browser.newPage();
  
  await page.setViewport({width: 1080, height: 1024});

  const notes = await getNotes(page);

  console.log(notes);

  await browser.close();
}

run();