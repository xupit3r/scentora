import Nano  from 'nano';
import csv from 'csvtojson';

export async function importFragrances () {
  const nano = Nano('http://admin:admin@127.0.0.1:5984');
  const scentsDb = nano.use('scents');
  const data = await csv({
    delimiter: ";"
  }).fromFile('./data/frangrances.csv');

  scentsDb.bulk({
    docs: data.map(doc => {
      return {
        id: doc.Perfume,
        url: doc.url,
        name: doc.Perfume.split('-').join(' '),
        notes: {
          top: doc.Top.split(',').map(s => s.trim()),
          middle: doc.Middle.split(',').map(s => s.trim()),
          base: doc.Base.split(',').map(s => s.trim())
        },
        accords: [
          doc.mainaccord1,
          doc.mainaccord2,
          doc.mainaccord3,
          doc.mainaccord4,
          doc.mainaccord5
        ].filter(a => a.trim().length),
        year: doc.Year,
        perfumer: doc.Perfumer1,
        country: doc.Country,
        gender: doc.gender
      };
    })
  });
}