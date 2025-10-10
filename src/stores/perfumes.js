import { reactive } from 'vue'
import { defineStore } from 'pinia'

const example = {
  "_id": "111e6b697952c8a5c6cdc3bbb6003cca",
  "_rev": "1-10ee23a4fbc6035825488a0bd8d97b53",
  "id": "rima-xi",
  "url": "https://www.fragrantica.com/perfume/carner-barcelona/rima-xi-17453.html",
  "name": "rima xi",
  "notes": {
    "top": [
      "cardamom",
      "saffron",
      "black pepper",
      "mint"
    ],
    "middle": [
      "ceylon cinnamon",
      "indonesian nutmeg",
      "coriander",
      "indian jasmine"
    ],
    "base": [
      "madagascar vanilla",
      "australian sandalwood",
      "benzoin",
      "amber",
      "musk",
      "virginian cedar"
    ]
  },
  "accords": [
    "warm spicy",
    "fresh spicy",
    "vanilla",
    "aromatic",
    "woody"
  ],
  "year": "2013",
  "perfumer": "sonia constant",
  "country": "Spain"
};

export const usePerfumeStore = defineStore('counter', () => {
  const perfume = reactive(example);

  return { perfume }
})