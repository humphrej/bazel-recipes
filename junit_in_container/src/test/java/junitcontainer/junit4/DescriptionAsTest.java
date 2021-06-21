package junitcontainer.junit4;

import junit.framework.Test;
import junit.framework.TestResult;
import org.junit.runner.Description;

/** Wraps {@link Description} into {@link Test} enough to fake . */
public class DescriptionAsTest implements Test {

  private final Description description;

  public DescriptionAsTest(Description description) {
    this.description = description;
  }

  public int countTestCases() {
    return 1;
  }

  public void run(TestResult result) {
    throw new UnsupportedOperationException();
  }

  /** determines the test name by reflection. */
  public String getName() {
    return description.getDisplayName();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }

    DescriptionAsTest that = (DescriptionAsTest) o;

    if (!description.equals(that.description)) {
      return false;
    }

    return true;
  }

  @Override
  public int hashCode() {
    return description.hashCode();
  }
}

