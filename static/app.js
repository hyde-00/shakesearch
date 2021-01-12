const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];

    if (results["count"] > 0) {
      rows.push(`<div style="font-weight:bold">Total: ${results["count"]} matches found.</div>`);
      rows.push(`<div> <div/>`);
      Object.keys(results["data"]).forEach(function(key) {
        values = results["data"][key];
        rows.push(`<div style="font-weight:bold">${key}:</div>`);
        for (let value of values) {
          if (value.length > 1) {
            rows.push(`<div>${value}</div>`);
          }
        }
        rows.push(`<div> <div/>`);
        rows.push(`<div> <div/>`);
      });
    } else {
      rows.push(`<div>No matches found.</div>`);
    }

    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
