import axios from "axios";
import { writeFileSync } from "fs";
const URL = "https://bitmidi.com/api/midi/random";
let data: { name: string; url: string }[] = [];
console.log("BitScraper started");

interface Manifest {
  name: string;
  author: string;
  software: string;
  source: string;
  manifestVersion: number;
  version: string;
  files: { name: string; url: string }[];
}
const MAX = 250;
async function scrap() {
  for (let i = 0; i < MAX; i++) {
    const res = await axios.get(URL);
    const result = res.data.result.result;
    const name: string = result.name;
    const url: string = result.downloadUrl;
    data.push({ name, url });
    console.log(
      "[" + (i + 1) + "/" + MAX + "]" + " Added " + result.name + " to memory"
    );
  }
}

function createManifest() {
  let manifest: Manifest = {
    author: "",
    name: "",
    source: "",
    software: "",
    manifestVersion: 0,
    version: "",
    files: [],
  };
  manifest.name = "250+ Midi Pack";
  manifest.author = "TheFakeWater";
  manifest.manifestVersion = 1;
  manifest.version = "1.0.0";
  manifest.source = "BitMidi";
  manifest.software = "BitScraper";

  for (let i = 0; i < data.length; i++) {
    const name = data[i].name;
    const url = data[i].url;
    manifest.files.push({ name, url });
  }
  writeFileSync("manifest.json", JSON.stringify(manifest));
  console.log("Created Manifest");
}
scrap().then(createManifest);
