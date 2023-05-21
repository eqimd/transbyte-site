import { newEditor } from "./editor.js";

var editorLeft, editorRight;

function openEditor(divName, isFocused, setter, text) {
  newEditor({containerId: divName, workerUrl: "worker.js"}).then(
    result => {
      setter(result);
      if (isFocused) result.focus();
      if (text) result.setText(text);
      console.log(`Editor started on ${divName}: ${result.saySomething()}`);
    }, error => console.error(error));
}


const textLeft = `class ClassOne {
    public static int funcOne() {
        return 1;
    }
}`

const textRight = `class ClassTwo {
    public static int funcTwo() {
        return 2;
    }
}`



//    var editorLeft = ace.edit("editorLeft");
//    editorLeft.setTheme("ace/theme/eclipse");
//    editorLeft.session.setMode("ace/mode/java");
//    editorLeft.setFontSize("13pt");

openEditor("editorLeft", true, a => editorLeft = a, textLeft);
// editor1.changeFont("JetBrainsMono-Regular", 13);

// console.info(editor1.getText());

//    var editorRight = ace.edit("editorRight");
//    editorRight.setTheme("ace/theme/eclipse");
//    editorRight.session.setMode("ace/mode/java");
//    editorRight.setFontSize("13pt");


openEditor("editorRight", false, a => editorRight = a, textRight);

    
const checkEquivButton = document.querySelector(".button_check_equiv");
const textAreaOutput = document.querySelector(".textarea_output");

// const xmlHttp = new XMLHttpRequest();
const URL = "http://localhost:534/";

// textAreaOutput.value = 'asdfasdfasfasdf'

function onClick() {
  // console.info("text1 = ");
  // console.info(editorLeft.getText());
  // editor1.IC("JetBrainsMono-Regular", 30)

  // console.info("text2 = ");
  // console.info(editorRight.getText());

  textAreaOutput.style.display = "none";

  checkEquivButton.classList.add("button_check_equiv_loading");

  // xmlHttp.open("GET", URL, true); // true for asynchronous 
  // xmlHttp.send(null);

  fetch(URL, {
    method: "POST",
    body: JSON.stringify({
      FirstClassCode: editorLeft.getText(),
      SecondClassCode: editorRight.getText(),
    })
    // mode: "no-cors",
  })
    .then((resp) => resp.json())
    .then((content) => {
      checkEquivButton.classList.remove("button_check_equiv_loading");
      textAreaOutput.value = content["Output"]
      textAreaOutput.style.display = "block";
    })

}

// function requestCallback(content) {
//   checkEquivButton.classList.remove("button_check_equiv_loading");
//   textAreaOutput.value = content;
//   textAreaOutput.style.display = "block";
// }

checkEquivButton.addEventListener("click", onClick);

// xmlHttp.onreadystatechange = function() { 
//     // if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
//     requestCallback(xmlHttp.responseText);
// }
