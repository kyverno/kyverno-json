import { AceEditor } from "./editor.js";

const selectInstance = NiceSelect.bind(document.getElementById("examples"));

const policyEditor = new AceEditor("policy-input");
const payloadEditor = new AceEditor("payload-input");

async function run() {
  const policy = policyEditor.getValue();
  const payload = payloadEditor.getValue();
  const output = document.getElementById("output");
  output.value = "Evaluating...";
  try {
    const reponse = await fetch("/api/playground/scan", {
      method: "POST",
      body: JSON.stringify({
        payload: payload,
        policy: policy,
      })
    })
    output.value = JSON.stringify(await reponse.json(), null, 2);
    output.style.color = "white";
  } catch (error) {
    output.value = error;
    console.error("Error:", error);
    output.style.color = "red";
  }
}

function share() {
  const payload = payloadEditor.getValue();
  const policy = policyEditor.getValue();

  const obj = {
    policy: policy,
    payload: payload,
  };

  const str = JSON.stringify(obj);
  var compressed_uint8array = pako.gzip(str);
  var b64encoded_string = btoa(
    String.fromCharCode.apply(null, compressed_uint8array)
  );

  const url = new URL(window.location.href);
  url.searchParams.set("content", b64encoded_string);
  window.history.pushState({}, "", url.toString());

  document.querySelector(".share-url__container").style.display = "flex";
  document.querySelector(".share-url__input").value = url.toString();
}

var urlParams = new URLSearchParams(window.location.search);
if (urlParams.has("content")) {
  const content = urlParams.get("content");
  try {
    const decodedUint8Array = new Uint8Array(
      atob(content)
        .split("")
        .map(function (char) {
          return char.charCodeAt(0);
        })
    );

    const decompressedData = pako.ungzip(decodedUint8Array, { to: "string" });
    if (!decompressedData) {
      throw new Error("Invalid content parameter");
    }
    const obj = JSON.parse(decompressedData);
    payloadEditor.setValue(obj.payload, -1);
    policyEditor.setValue(obj.policy, -1);
  } catch (error) {
    console.error(error);
  }
}

function copy() {
  const copyText = document.querySelector(".share-url__input");
  copyText.select();
  copyText.setSelectionRange(0, 99999);
  navigator.clipboard.writeText(copyText.value);
  window.getSelection().removeAllRanges();

  const tooltip = document.querySelector(".share-url__tooltip");
  tooltip.style.opacity = 1;
  setTimeout(() => {
    tooltip.style.opacity = 0;
  }, 3000);
}

const runButton = document.getElementById("run");
const shareButton = document.getElementById("share");
const copyButton = document.getElementById("copy");

runButton.addEventListener("click", run);
shareButton.addEventListener("click", share);
copyButton.addEventListener("click", copy);
document.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.code === "Enter") {
    run();
  }
});

fetch("./assets/data.json")
  .then((response) => response.json())
  .then(({ examples }) => {

    // Load the examples into the select element
    const examplesList = document.getElementById("examples");
    examples.forEach((example) => {
      const option = document.createElement("option");
      option.value = example.name;
      option.innerText = example.name;

      if (example.name === "default") {
        if (!urlParams.has("content")) {
          payloadEditor.setValue(example.payload, -1);
          policyEditor.setValue(example.policy, -1);
        }
      } else {
        examplesList.appendChild(option);
      }
    });

    selectInstance.update();

    examplesList.addEventListener("change", (event) => {
      const example = examples.find(
        (example) => example.name === event.target.value
      );
      payloadEditor.setValue(example.payload, -1);
      policyEditor.setValue(example.policy, -1);
    });
  })
  .catch((err) => {
    console.error(err);
  });