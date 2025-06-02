interface StreamElementsEvent {
    detail: {
        fieldData: Record<string, string>;
    };
}

declare global {
    interface WindowEventMap {
        onEventReceived: StreamElementsEvent;
        onWidgetLoad: StreamElementsEvent;
    }
}
