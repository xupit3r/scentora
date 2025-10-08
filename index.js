import puppeteer from 'puppeteer';
import Nano  from 'nano';

const nano = Nano('http://admin:admin@127.0.0.1:5984');


async function getFragranceNotes (page) {
  const TOP_NOTES = 'Top Notes'
  const MIDDLE_NOTES = 'Middle Notes';
  const BASE_NOTES = 'Base Notes';

  const text = await page.$eval('#pyramid', result => {
    return result.innerText
  });
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

async function getPageInfo (page, url) {
  await page.goto(url);

  await delay(3000);

  const notes = await getFragranceNotes(page);
  const accords = await getFrangranceAccords(page);
  const description = await getFrangranceDescription(page);

  return {
    notes,
    accords,
    description
  }
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
  const notesDb = nano.use('notes');

  await page.goto("https://www.fragrantica.com/notes/");

  const notes = await page.$$eval('.notebox a', results => {
    return results.map(el => {
      return {
        note: el.innerText.trim(),
        url: el.href
      }
    });
  });

  await notesDb.bulk({
    docs: notes
  });

  return notes;
}

async function getScentsByNote (page, url) {
  await page.goto(url);

  const scents = await page.$$eval('.prefumeHbox a', results => {
    return results.map(el => {
      return {
        name: el.innerText.trim(),
        url: el.href
      };
    });
  });

  const infos = []
  for (let i = 0; i < scents.length; i++) {
    const info = await getPageInfo(page, scents[i].url);
    info.url = scent[i].url;
    info.name = scent[i].name;

    infos.push(info);
  }

  return infos;
}

async function readNotes () {
  const notesDb = nano.use('notes');
  const results = await notesDb.list({ include_docs: true });

  return results.rows.map(({ doc }) => doc);
}

async function run () {
  const browser = await puppeteer.launch({
    headless: false
  });

  const page = await browser.newPage();
  await page.setViewport({width: 1080, height: 1024});

  const scents = await getScentsByNote(
    page,
    'https://www.fragrantica.com/notes/Apple-146.html'
  );

  console.log(scents);

  await browser.close();
}

run();