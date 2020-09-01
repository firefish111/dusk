const getlib = require("node-fetch");
const fs = require("fs");
arg = process.argv.slice(2);
console.log(arg);
(async () => { 
  await fs.mkdir("dusk_modules", err => err);
  if (arg.length == 0) console.log(fs.readFileSync("help"))
  else {
    console.log("DB0");
    arg.forEach(async pkg => {
      console.log("DB1");
      let ver = [];
      let code;
      await getlib(`https://duskcdn.firefish.repl.co/cdn/${pkg}`, { method: "GET", headers: { //
        "X-Requested-With": "night-dusk-pkg" // hahaNO
      } })
        .then(res => res.text())
        .then(text => code = text)
        //.then(res => ver = res.headers["X-Package-Version"]);
        console.log(ver);
      fs.writeFileSync(`${__dirname}/dusk_modules/${pkg}.night`, code)
    });
  }
})();