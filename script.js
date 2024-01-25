// The DApp listens to announced providers
window.addEventListener("eip6963:announceProvider", (event) => {
  console.log(event);
});

// The DApp dispatches a request event which will be heard by
// Wallets' code that had run earlier
window.dispatchEvent(new Event("eip6963:requestProvider"));
