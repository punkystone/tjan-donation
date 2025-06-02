window.addEventListener("onWidgetLoad", function (data) {
    const format = data.detail.fieldData["format"];
    const months = [
        "Tjanuar",
        "Tjebruar",
        "Tjärz",
        "Tjapril",
        "Tjai",
        "Tjuni",
        "Tjuli",
        "Tjaugust",
        "Tjeptember",
        "Tjoktober",
        "Tjovember",
        "Tjecember",
    ];
    const monthElement = document.querySelector(".month");
    const donationElement = document.querySelector(".donation");
    if (donationElement === null) return;
    if (monthElement === null) return;
    if (format === undefined) return;
    const donationUrl = "wss://tjan.punkystone.de/donation";
    startSocket(donationUrl, donationElement);
    const updateLoop = (): void => {
        const now = new Date();
        const monthName = months[now.getMonth()];
        if (monthName === undefined) return;
        monthElement.textContent = format.replace("{month}", monthName);
    };
    updateLoop();
    setInterval(updateLoop, 1000);
});

const startSocket = (donationUrl: string, donationElement: Element): void => {
    const wss = new WebSocket(donationUrl);
    wss.addEventListener("message", (event: MessageEvent): void => {
        const data = event.data as string;
        donationElement.innerHTML = data;
    });
    wss.addEventListener("close", (): void => {
        console.log("Connection closed, retrying in 5 seconds");
        setTimeout(function () {
            startSocket(donationUrl, donationElement);
        }, 5000);
    });
};
