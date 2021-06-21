package junitcontainer;

/** Represents a class under test */
public class UnderTest {

  /**
   * Returns the magic number
   *
   * @return the magic number
   */
  public static int magicNumber() {
    return "open sesame".hashCode();
  }
}

