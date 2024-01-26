import { EIP6963AnnounceProviderEvent, EIP6963ProviderDetail } from "./types";

export function registerEIP6963(providers: EIP6963ProviderDetail[]) {
  window.addEventListener(
    "eip6963:announceProvider",
    (event: EIP6963AnnounceProviderEvent) => {
      providers.push(event.detail);
    },
  );

  window.dispatchEvent(new Event("eip6963:requestProvider"));
}
