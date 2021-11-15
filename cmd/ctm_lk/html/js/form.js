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

function zero_first_format(value) {
    if (value < 10) {
        value = '0' + value;
    }
    return value
}
/* функция получения текущей даты и времени */
function date_() {
    var day = zero_first_format(current_datetime.getDate());
    var month = zero_first_format(current_datetime.getMonth() + 1);
    var year = current_datetime.getFullYear();

    //return day+"."+month+"."+year+" "+hours+":"+minutes+":"+seconds;
    return year + "-" + month + "-" + day
}
function time_() {
    var hours = zero_first_format(current_datetime.getHours());
    var minutes = zero_first_format(current_datetime.getMinutes());
    var seconds = zero_first_format(current_datetime.getSeconds());

    return hours + ":" + minutes
}

function defaultValue() {
    document.getElementById('datein').value = date_();
    document.getElementById('timein').value = time_();
}

defaultValue()

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
    var response = await fetch(evt.target.action, {
        method: 'POST',
        body: new FormData(form)
    });
    if (response.ok) {
        onFormPostSuccess()
    } else {
        onFormPostError()
    }
}