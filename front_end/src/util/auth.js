const isClient = typeof window !== undefined;

export function setToken(data) {
  try {
    if (isClient) {
      window.localStorage.setItem("AccessToken", data);
    } else {
      throw new Error(`isClient: ${isClient}`);
    }
  } catch (error) {
    console.warn("setToken", error);
  }
}

export function getToken() {
  try {
    if (isClient) {
      return window.localStorage.getItem("AccessToken")
    } else {
      throw new Error(`isClient: ${isClient}`);
      
    }
  } catch (error) {
    console.warn("getToken", error);
    return{status:false}
  }
}
