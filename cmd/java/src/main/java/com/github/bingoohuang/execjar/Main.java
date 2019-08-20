package com.github.bingoohuang.execjar;

import lombok.SneakyThrows;

import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

/**
 * Just a main class.
 */
public class Main {

    /**
     * A main method of this program which finds two classes and prints their names
     * from guava and guice, respectively.
     *
     * @param args command line arguments.
     */
    @SneakyThrows
    public static void main(String[] args) {
        System.out.println(Class.forName("com.google.common.io.BaseEncoding"));

        List<Double> list = new ArrayList<>();

        for (int i = 0; i < 9000000; i++) {
            list.add(Math.random());
        }

        String inData;
        Scanner scan = new Scanner(System.in);

        FOR:
        for (; ; ) {
            inData = scan.nextLine();
            switch (inData) {
                case "quit":
                    break FOR;
                case "sleep":
                    Thread.sleep(10000);
                    System.out.println("You entered:" + inData);
                    break;
                default:
                    System.out.println("You entered:" + inData);
            }
        }

        scan.close();
    }
}
