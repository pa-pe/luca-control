const anim_time = 500;
const iconSuccess = '<i class="bi-check-circle" style="color: green;"></i>';
const iconLoading = '<i class="bi-arrow-repeat" style="color: blue;"></i>';
const iconFail = '<i class="bi-x-circle" style="color: red;"></i>';

function send_update_data(data){
//    alert(data);
//console.log(data);
//return "";
    fetch('/update_model', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
//        body: JSON.stringify({ "model": model, "id": ID, [field]: value }),
        body: JSON.stringify(data),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            query_result = 1;
//            statusContainer.animShow(iconSuccess);
alert("update ok");
        } else {
            query_result = 0;
//            statusContainer.animShow(iconFail);
            alert('Failed to update: ' + (data.error || 'Unknown error'));
        }
    })
    .catch(error => {
        query_result = 0;
//        statusContainer.animShow(iconFail);
        alert('Error: ' + error);
    });
}

function showQueuedAnim(name, tagCont, doFunc) {
    if ($(tagCont).hasClass("change_animation_progress")) {
        setTimeout(showQueuedAnim, 100, name, tagCont, doFunc);
        return;
    }
    doFunc();
}

function animHideElement(element, animTime, tagCont, callback) {
    $(tagCont).addClass("change_animation_progress");
    $(element).animate({ opacity: 0 }, animTime, function() {
        element.style.display = "none";
        $(tagCont).removeClass("change_animation_progress");
        if (callback) callback();
    });
}

function animShowElement(element, animTime, tagCont) {
    showQueuedAnim("animShow", tagCont, function() {
        element.style.display = "block";
        $(tagCont).addClass("change_animation_progress");
        $(element).animate({ opacity: 1 }, animTime, function() {
            $(tagCont).removeClass("change_animation_progress");
        });
    });
}

function animHideStatusContainer(statusContainer, animTime, tagCont) {
    if (statusContainer.style.display != "none") {
        showQueuedAnim("hide-cont", tagCont, function() {
            $(tagCont).addClass("change_animation_progress");
            $(statusContainer).animate({ opacity: 0 }, animTime, function() {
                statusContainer.style.display = "none";
                statusContainer.innerHTML = '';
                $(tagCont).removeClass("change_animation_progress");
            });
        });
    }
}

function animShowStatusContainer(statusContainer, iconTag, animTime, tagCont, dropdown) {
    animHideStatusContainer(statusContainer, animTime, tagCont);

    showQueuedAnim("show-success", tagCont, function() {
        $(tagCont).addClass("change_animation_progress");
        statusContainer.style.opacity = "0";
        statusContainer.innerHTML = iconTag;
        statusContainer.style.display = "block";
        $(statusContainer).animate({ opacity: 1 }, animTime, function() {
            $(tagCont).removeClass("change_animation_progress");

            if (iconTag == iconSuccess) {
                animHideStatusContainer(statusContainer, animTime, tagCont);
                animShowElement(dropdown, animTime, tagCont);
            }
        });
    });
}
