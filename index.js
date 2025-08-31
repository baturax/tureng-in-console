import readline from "readline";
import fetch from "node-fetch";
import { JSDOM } from "jsdom";

async function getTerm() {
  const argTerm = process.argv[2];
  if (argTerm) return argTerm;

  // Eğer argüman yoksa kullanıcıdan input al
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  });

  return new Promise(resolve => {
    rl.question("Enter a word: ", (answer) => {
      rl.close();
      resolve(answer.trim());
    });
  });
}

async function fetchTureng(term) {
  const url = `https://tureng.com/en/turkish-english/${term}`;

  const res = await fetch(url, {
    headers: { "User-Agent": "Mozilla/5.0" },
  });

  const html = await res.text();
  const dom = new JSDOM(html);
  const document = dom.window.document;

  const results = [];
  const rows = document.querySelectorAll("tr.tureng-manual-stripe-even, tr.tureng-manual-stripe-odd");

  rows.forEach((row, i) => {
    const english = row.querySelector("td.en.tm a")?.textContent.trim();
    const turkish = row.querySelector("td.tr.ts a")?.textContent.trim();
    if (english && turkish) results.push(`${i + 1}: ${english} -> ${turkish}`);
  });

  return results;
}

const term = await getTerm();
if (!term) {
  console.log("No word provided!");
  process.exit(1);
}

const results = await fetchTureng(term);
console.log(results.join("\n"));

