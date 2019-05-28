const app = {};
const init = () => {
  console.log("App is running.");
  app.data = {};
  addListeners();
};

const getValues = (e) => {
    e.preventDefault();
    var formData = new FormData(e.target);
    formData.forEach((value, key) => {
      app.data[key] = value.replace(/["']/g, "");
    });
    data = cleanValue(app.data);
    console.log(data);
    send(data)
};

function send(data) {
  fetch("https://templategen.com/", {
    method: "post",
    body: JSON.stringify(data)
  }).then(function(response) {
    console.log(response)
    download(data.templateGroupName);
  }).catch((err) => {
    console.log(err)
    document.querySelector('.error').style.display = "block"
  })
}

function download(zipName){
  document.querySelector('.error').style.display = "none"
  var url= "https://templategen.com/download/" + zipName;
  location.assign(url);
}

const addListeners = () => {
  console.log("Adding event listeners.");

  var hide = document.getElementsByClassName("hide");
  var txt = document.getElementsByClassName("text-input");

  for (let i = 0; i < hide.length; i++) {
    resize(hide[i], txt[i]);
    txt[i].addEventListener("input", () => resize(hide[i], txt[i]));
  }

  function resize(hide, txt) {
    hide.textContent = txt.value;
    txt.style.width = hide.clientWidth + "px";
  }

  document.querySelector("form").addEventListener("submit", getValues);
};


function cleanValue (object) {
  // clear whitespace and split by comma into an array
  object.sizes = object.sizes.replace(/\s+/g, "");
  object.sizes = object.sizes.split(",");
  object.start = object.start.replace(/\s+/g, "");
  object.start = object.start.split(",");
  object.middle = object.middle.replace(/\s+/g, "");
  object.middle = object.middle.split(",");
  object.end = object.end.replace(/\s+/g, "");
  object.end = object.end.split(",");

  //str to int
  object.frameLimit = Number(object.frameLimit);
  object.frameMinCount = Number(object.frameMinCount);
  object.baseSize = Number(object.baseSize);

  //todo
  object.templateGroupName = object.templateGroupName.replace(/\s+/g, "-");
  object.templateName = object.templateName.replace(/\s+/g, "-");
  object.templateSet = object.templateName;
  return object;
}

window.onload = init;
