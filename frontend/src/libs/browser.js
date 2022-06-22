export class BooksysBrowser {
  static isMobileResponsive(maxWidth=992) {
    const screenWidth = window.innerWidth;
    if(maxWidth < screenWidth){
      return false;
    }
    return true;
  }
}
