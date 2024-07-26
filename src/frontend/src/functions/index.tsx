
export function capitalize(str: string) {
    str = str[0].toUpperCase() + str.slice(1)
    return str
}

export function subtitleError(msg: string) {
    var title = document.getElementById("undertitle");
    if (title)
        title.innerText = msg;
}
