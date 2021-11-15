/* функция добавления ведущих нулей */
/* (если число меньше десяти, перед числом добавляем ноль) */
var current_datetime = new Date();
var popup = document.querySelector(".frm_posted");
var close = popup.querySelectorAll(".modal-close-button");
for (var i = 0; i < close.length; i++) {
    close[i].addEventListener("click", function (evt) {
        evt.preventDefault();
        popup.classList.remove("modal-show");
        /*window.location.reload();*/
    });
}

/*отправка формы*/
var form = document.querySelector('.form_main');

function onFormPostError() {
    var header = popup.querySelector('h2')
    header.value = "Ошибка отправки сообщения"
    popup.classList.add('modal-show');
};

function onFormPostSuccess() {
    popup.classList.add('modal-show');
    form.reset()
    document.getElementById('datein').value = date_();
    document.getElementById('timein').value = time_();
};

form.addEventListener('submit', function (evt) {
    submit(evt)
    evt.preventDefault();
});

async function submit(evt) {
    var object = {};
    var formData = new FormData(form);
    formData.forEach(function (value, key) {
        object[key] = value;
    });
    var json = JSON.stringify(object);
    var response = await fetch(evt.target.action, {
        headers: {
            // 'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: 'POST',
        body: json
    });
    if (response.ok) {
        onFormPostSuccess()
    } else {
        onFormPostError()
    }
}