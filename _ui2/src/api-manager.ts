export interface XhrOptions {
  method?: string;
  params?: any;
}

export class ApiManager {
  public static async xhrJson(url: string, options?: XhrOptions) {
    const response = await this.xhrText(url, options);
    return JSON.parse(response || 'null');
  }

  public static async xhrText(url: string, options?: XhrOptions) {
    const request = await this.xhr(url, options);
    return (request as any).responseText;
  }

  public static xhr(url: string, options?: XhrOptions) {
    const request = new XMLHttpRequest();
    return new Promise((resolve, reject) => {
      request.onreadystatechange = () => {
        if (request.readyState === 4) {
          if (request.status === 200) {
            try {
              resolve(request);
            } catch (e) {
              reject(e);
            }
          } else {
            reject(request);
          }
        }
      };
      const method = (options && options.method) || "GET";
      request.open(method, url);
      const params = (options && options.params) || {};
      if (params) {
        request.setRequestHeader("Content-Type", "application/json");
        request.send(JSON.stringify(params));
      }
    })
  }
}
