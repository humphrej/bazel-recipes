package junitcontainer.junit4;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;

@RunWith(Suite.class)
@Suite.SuiteClasses({
  MagicNumberTest.class,
})
public class TestSuite {}

